let card;

document.addEventListener('DOMContentLoaded', () => {
    const stripe = initializeStripe();
    setupFormListeners(stripe);
});

const initializeStripe = () => {
    const pkey = document.getElementById('stripe_public_key')?.innerHTML;
    if (!pkey) {
        showCardError('Failed to load payment processor. Please try again.');
        throw new Error('Stripe public key missing');
    }
    return Stripe(pkey);
};

const setupFormListeners = (stripe) => {
    const form = document.getElementById('charge_form');
    const amountInput = document.getElementById('amount');

    setupInputValidation(amountInput);

    form.addEventListener('submit', async (event) => {
        event.preventDefault();
        const amount = amountInput.value.trim();

        if (!validateAmountInput(amount)) {
            return;
        }

        toggleProcessingState(form, true);
        try {
            const payload = {
                amount: amount,
                currency: 'brl',
            };

            const clientSecret = await fetchPaymentIntent(payload);
            await confirmPayment(stripe, clientSecret, form);
        } catch (err) {
            showCardError(err);
        } finally {
            toggleProcessingState(form, false);
        }
    });

    setupCardElements(stripe);
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

const validateAmountInput = (value) => {
    const amount = parseFloat(value);
    if (isNaN(amount) || amount <= 0) {
        showCardError('Please enter a valid amount.');
        return false;
    }
    return true;
};

const fetchPaymentIntent = async (payload) => {
    const apiUrl = document.getElementById('api_url')?.innerHTML;
    if (!apiUrl) {
        showCardError('Failed to load API url. Please try again.');
    }

    try {
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
    } catch (err) {
        console.error('Error fetching payment intent:', err);
        showCardError('Failed to create payment intent.');
        throw err;
    }
};

const confirmPayment = async (stripe, clientSecret, form) => {
    const cardholderName = document.getElementById('cardholder-name').value;

    try {
        const result = await stripe.confirmCardPayment(clientSecret, {
            payment_method: {
                card,
                billing_details: {
                    name: cardholderName,
                },
            },
        });

        if (result.error) {
            throw new Error(result.error.message);
        } else if (result.paymentIntent && result.paymentIntent.status === 'succeeded') {
            processPaymentSuccess(result.paymentIntent);
            setTimeout(() => form.submit(), 1000);
        }
    } catch (err) {
        showCardError(err.message);
        throw err;
    }
};

const processPaymentSuccess = (paymentIntent) => {
    document.getElementById('payment_method').value = paymentIntent.payment_method_types[0];
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
    console.error('Error occurred:', message);
    cardMessages.classList.add('alert-danger');
    cardMessages.classList.remove('alert-success', 'd-none');
    cardMessages.innerText = message || 'An unexpected error occurred.';
};

const showCardSuccess = () => {
    const cardMessages = document.getElementById('card-messages');
    cardMessages.classList.add('alert-success');
    cardMessages.classList.remove('alert-danger', 'd-none');
    cardMessages.innerText = 'Transaction successful!';
};
