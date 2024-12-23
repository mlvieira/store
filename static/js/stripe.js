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

        toggleProcessingState(form, true);
        try {
            if (paymentMode === 'subscription') {
                await handleSubscription(stripe, email, form);
            } else {
                await handleOneTimePayment(stripe, amountInput, form);
            }
        } catch (err) {
            console.error(err);
            showCardError(err.message);
        } finally {
            toggleProcessingState(form, false);
        }
    });

    setupCardElements(stripe);
};

const handleSubscription = async (stripe, email, form) => {
    const planId = document.getElementById('plan_id')?.value;
    if (!planId) {
        throw new Error('Plan ID is required for subscription.');
    }

    const paymentMethod = await createPaymentMethod(stripe, email);

    const clientSecret = await createSubscription(paymentMethod, email, planId);

    await confirmPayment(stripe, clientSecret, form);
};

const handleOneTimePayment = async (stripe, amountInput, form) => {
    if (!validateAmountInput(amountInput)) {
        throw new Error('Invalid amount.');
    }

    const amountInCents = Math.round(parseFloat(amountInput) * 100);

    const clientSecret = await createPaymentIntent(amountInCents);

    await confirmPayment(stripe, clientSecret, form);
};

const createPaymentMethod = async (stripe, email) => {
    const cardholderName = document.getElementById('cardholder-name').value;

    const { paymentMethod, error } = await stripe.createPaymentMethod({
        type: 'card',
        card: card,
        billing_details: {
            name: cardholderName,
            email: email,
        },
    });

    if (error) {
        throw new Error(error.message);
    }
    return paymentMethod;
};

const createSubscription = async (paymentMethod, email, planId) => {
    const payload = {
        email: email,
        plan_id: planId,
        payment_method: paymentMethod.Id,
        last_four: paymentMethod.last4,
    };

    const response = await fetch(`${apiUrl}/api/create-subscription`, {
        method: 'POST',
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload),
    });

    if (!response.ok) {
        throw new Error(`HTTP Error: ${response.status}`);
    }

    const data = await response.json();
    if (!data.client_secret) {
        throw new Error('Missing client_secret in response');
    }

    return data.client_secret;
};

const createPaymentIntent = async (amount) => {
    const payload = {
        amount: amount,
        currency: 'brl',
    };

    const response = await fetch(`${apiUrl}/api/payment-intent`, {
        method: 'POST',
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload),
    });

    if (!response.ok) {
        throw new Error(`HTTP Error: ${response.status}`);
    }

    const data = await response.json();
    if (!data.client_secret) {
        throw new Error('Missing client_secret in response');
    }

    return data.client_secret;
};

const confirmPayment = async (stripe, clientSecret, form) => {
    try {
        const result = await stripe.confirmCardPayment(clientSecret);
        if (result.error) {
            throw new Error(result.error.message);
        } else if (result.paymentIntent && result.paymentIntent.status === 'succeeded') {
            processPaymentSuccess(result.paymentIntent);
            setTimeout(() => form.submit(), 1000);
        }
    } catch (err) {
        console.error(err);
        showCardError(err.message);
        throw err;
    }
};

const processPaymentSuccess = (paymentIntent) => {
    document.getElementById('payment_method').value = paymentIntent.payment_method;
    document.getElementById('payment_intent').value = paymentIntent.id;
    document.getElementById('payment_amount').value = paymentIntent.amount;
    document.getElementById('payment_currency').value = paymentIntent.currency;

    showCardSuccess();
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

    payButton.disabled = isProcessing;
    processing.classList.toggle('d-none', !isProcessing);
    form.querySelectorAll('input, button').forEach((element) => {
        element.disabled = isProcessing;
    });
};

const showCardError = (message) => {
    const cardMessages = document.getElementById('card-messages');
    cardMessages.classList.add('alert-danger');
    cardMessages.classList.remove('alert-success', 'd-none');
    cardMessages.innerText = message || 'An unexpected error occurred.';
};

const showCardSuccess = () => {
    const cardMessages = document.getElementById('card-messages');
    cardMessages.classList.add('alert-success');
    cardMessages.classList.remove('alert-danger', 'd-none');
    cardMessages.innerText = 'Payment successful!';
};
