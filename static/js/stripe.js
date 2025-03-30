let card;
let apiUrl = '';

document.addEventListener('DOMContentLoaded', () => {
    const stripe = initializeStripe();
    initGlobalConfig();
    setupFormListeners(stripe);
});

const initializeStripe = () => {
    const pkey = document.getElementById('stripe_public_key')?.innerText;
    if (!pkey) {
        showCardError('Failed to load payment processor. Please try again.');
        throw new Error('Stripe public key missing');
    }
    return Stripe(pkey);
};

const initGlobalConfig = () => {
    apiUrl = document.getElementById('api_url')?.innerText;
    if (!apiUrl) {
        showCardError('Failed to load API URL. Please try again.');
        throw new Error('API URL missing');
    }
};

const setupFormListeners = (stripe) => {
    const form = document.getElementById('charge_form');
    const paymentMode = document.getElementById('payment_mode').value;
    const amountInput = document.getElementById('amount');

    if (paymentMode !== 'subscription') {
        setupInputValidation(amountInput);
    }

    form.addEventListener('submit', async (event) => {
        event.preventDefault();
        const email = document.getElementById('email').value.trim();

        if (!validateEmail(email)) {
            showCardError('Please provide a valid email.');
            return;
        }

        if (
            paymentMode !== 'subscription' &&
            !validateAmountInput(amountInput.value)
        ) {
            return;
        }

        toggleProcessingState(form, true);
        try {
            const paymentMethod = await createPaymentMethod(stripe, email);
            let clientSecret;

            if (paymentMode === 'subscription') {
                const planId = document.getElementById('plan_id')?.value;
                if (!planId) {
                    throw new Error('Plan ID is required for subscription.');
                }
                clientSecret = await createSubscription(
                    paymentMethod,
                    email,
                    planId,
                    amountInput
                );
            } else {
                const amountInCents = Math.round(
                    parseFloat(amountInput.value) * 100
                );
                clientSecret = await createPaymentIntent(
                    amountInCents,
                    paymentMethod.id
                );
            }

            await confirmIntent(
                stripe,
                clientSecret,
                form,
                paymentMethod.id,
                paymentMode
            );
        } catch (err) {
            console.error('Form submission error:', err);
            showCardError(
                err.message || 'An unexpected error occurred.'
            );
        } finally {
            toggleProcessingState(form, false);
        }
    });

    setupCardElements(stripe);
};

const createPaymentMethod = async (stripe, email) => {
    const { paymentMethod, error } = await stripe.createPaymentMethod({
        type: 'card',
        card: card,
        billing_details: {
            email: email,
        },
    });

    if (error) {
        console.error('createPaymentMethod error:', error);
        throw new Error(error.message);
    }
    return paymentMethod;
};

const createSubscription = async (
    paymentMethod,
    email,
    planId,
    amountInput
) => {
    const payload = {
        email: email,
        plan_id: planId,
        payment_method: paymentMethod.id,
        last_four: paymentMethod.card.last4,
        card_brand: paymentMethod.card.brand,
        expiry_month: paymentMethod.card.exp_month,
        expiry_year: paymentMethod.card.exp_year,
        first_name: document.querySelector('#first-name').value,
        last_name: document.querySelector('#last-name').value,
        product_id: document.querySelector('input[name="widget_id"]').value,
        currency: 'brl',
        amount: Math.round(parseFloat(amountInput.value) * 100),
    };

    console.log('Sending payload to create subscription:', payload);
    const response = await fetch(`${apiUrl}/api/create-subscription`, {
        method: 'POST',
        headers: {
            Accept: 'application/json',
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload),
    });

    console.log('Create Subscription Response Status:', response.status);
    if (!response.ok) {
        const errorText = await response.text();
        console.error('Create Subscription Error Response:', errorText);
        throw new Error(`HTTP Error: ${response.status} - ${errorText}`);
    }

    const data = await response.json();
    console.log('Create Subscription Response Data:', data);
    if (!data.content) {
        console.error(
            'Missing client_secret (expected in content field) in subscription response:',
            data
        );
        throw new Error('Missing client_secret for setup in response');
    }

    return data.content;
};

const createPaymentIntent = async (amount, paymentMethodId) => {
    const payload = {
        amount: amount,
        currency: 'brl',
        payment_method: paymentMethodId,
    };

    console.log('Sending payload to create PI:', payload);
    const response = await fetch(`${apiUrl}/api/payment-intent`, {
        method: 'POST',
        headers: {
            Accept: 'application/json',
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload),
    });

    console.log('Create PI Response Status:', response.status);
    if (!response.ok) {
        const errorText = await response.text();
        console.error('Create PI Error Response:', errorText);
        throw new Error(`HTTP Error: ${response.status} - ${errorText}`);
    }

    const data = await response.json();
    console.log('Create PI Response Data:', data);
    if (!data.client_secret) {
        console.error(
            'Missing client_secret in payment intent response:',
            data
        );
        throw new Error('Missing client_secret in response');
    }

    return data.client_secret;
};

const confirmIntent = async (
    stripe,
    clientSecret,
    form,
    paymentMethodId,
    paymentMode
) => {
    let confirmationPromise;

    if (paymentMode !== 'subscription') {
        console.log(
            'Attempting confirmCardPayment with clientSecret:',
            clientSecret,
            'and paymentMethodId:',
            paymentMethodId
        );
        confirmationPromise = stripe.confirmCardPayment(clientSecret, {
            payment_method: paymentMethodId,
        });
    } else {
        console.log(
            'Attempting confirmCardSetup with clientSecret:',
            clientSecret,
            'and paymentMethodId:',
            paymentMethodId
        );
        confirmationPromise = stripe.confirmCardSetup(clientSecret, {
            payment_method: paymentMethodId,
        });
    }

    const result = await confirmationPromise;
    console.log('Stripe confirmation result:', result);

    if (result.error) {
        console.error('Stripe confirmation error:', result.error);
        throw new Error(result.error.message);
    }

    let intent = null;
    if (
        result.paymentIntent &&
        result.paymentIntent.status === 'succeeded'
    ) {
        console.log('PaymentIntent Succeeded:', result.paymentIntent);
        intent = result.paymentIntent;
    } else if (
        result.setupIntent &&
        result.setupIntent.status === 'succeeded'
    ) {
        console.log('SetupIntent Succeeded:', result.setupIntent);
        intent = result.setupIntent;
    } else {
        const status =
            result.paymentIntent?.status ||
            result.setupIntent?.status ||
            'unknown';
        const intentType = result.paymentIntent ? 'Payment' : 'Setup';
        console.warn(`Unhandled ${intentType}Intent status: ${status}`, result);
        throw new Error(
            `Payment/Setup status: ${status}. Further action may be required.`
        );
    }

    processPaymentSuccess(intent);
    setTimeout(() => {
        console.log('Submitting form after success.');
        form.submit();
    }, 1000);
};

const processPaymentSuccess = (intent) => {
    console.log('Processing success for intent:', intent);
    showCardSuccess();

    const paymentMethodInput = document.getElementById('payment_method');
    if (paymentMethodInput && intent.payment_method) {
        paymentMethodInput.value =
            typeof intent.payment_method === 'string'
                ? intent.payment_method
                : intent.payment_method.id;
    }

    const intentIdInput = document.getElementById('intent_id');
    if (intentIdInput) {
        intentIdInput.value = intent.id;
    }

    const amountInput = document.getElementById('payment_amount');
    const currencyInput = document.getElementById('payment_currency');

    if (intent.object === 'payment_intent') {
        console.log('Processing PaymentIntent specific data');
        if (amountInput) amountInput.value = intent.amount;
        if (currencyInput) currencyInput.value = intent.currency;
    } else if (intent.object === 'setup_intent') {
        console.log('Processing SetupIntent specific data');
        if (amountInput) amountInput.value = ''; // Or 0
        if (currencyInput) currencyInput.value = '';
    }
};

const setupCardElements = (stripe) => {
    const elements = stripe.elements();
    const style = {
        base: {
            fontSize: '16px',
            lineHeight: '24px',
        },
    };

    card = elements.create('card', { style, hidePostalCode: true });
    card.mount('#card-element');

    card.addEventListener('change', (event) => {
        const displayError = document.getElementById('card-errors');
        if (event.error) {
            displayError.classList.remove('d-none');
            displayError.textContent = event.error.message;
        } else {
            displayError.classList.add('d-none');
            displayError.textContent = '';
        }
    });
};

const validateAmountInput = (value) => {
    const amount = parseFloat(value);
    if (isNaN(amount) || amount <= 0) {
        showCardError('Please enter a valid amount.');
        return false;
    }
    return true;
};

const validateEmail = (email) => {
    const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailPattern.test(email);
};

const setupInputValidation = (amountInput) => {
    amountInput.addEventListener('input', () => {
        amountInput.value = amountInput.value.replace(/[^0-9.]/g, '');
        const parts = amountInput.value.split('.');
        if (parts.length > 2) {
            amountInput.value = parts[0] + '.' + parts[1];
        }
    });

    amountInput.addEventListener('blur', () => {
        const value = parseFloat(amountInput.value);
        if (!isNaN(value)) {
            amountInput.value = value.toFixed(2);
        } else {
            amountInput.value = '';
        }
    });
};

const toggleProcessingState = (form, isProcessing) => {
    const payButton = document.getElementById('pay-button');
    const processing = document.getElementById('processing-payment');

    if (payButton) payButton.disabled = isProcessing;
    if (processing) processing.classList.toggle('d-none', !isProcessing);

    form.querySelectorAll('input, button').forEach((element) => {
        if (element.id !== 'card-element' && !element.closest('#card-element')) {
             element.disabled = isProcessing;
        }
    });
};

const showCardError = (message) => {
    const cardMessages = document.getElementById('card-messages');
    if (!cardMessages) return;
    cardMessages.classList.add('alert-danger');
    cardMessages.classList.remove('alert-success', 'd-none');
    cardMessages.innerText = message || 'An unexpected error occurred.';
};

const showCardSuccess = () => {
    const cardMessages = document.getElementById('card-messages');
    if (!cardMessages) return;
    cardMessages.classList.add('alert-success');
    cardMessages.classList.remove('alert-danger', 'd-none');
    cardMessages.innerText = 'Payment successful!';
};
