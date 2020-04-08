// begindoc: all
// https://dashboard.stripe.com/account/apikeys
// It is fine to publish the Stripe Publishable key,
// as it has no dangerous permissions.
const PUBLISHABLE_API_KEY = "pk_test_1234567890abcdef";

const client = require("stripe-client")(PUBLISHABLE_API_KEY);

interface StripeCard {
  address_zip: string;
  cvc?: string; // usually required
  exp_month: string; // two-digit number
  exp_year: string; // two- or four-digit number
  number: string;
}

const createCardToken = async (card: StripeCard): Promise<string | null> => {
  // https://stripe.com/docs/api/tokens/create_card
  const response = await client.createToken({ card });

  if (response && response.id) {
    return response.id;
  } else {
    return null;
  }
};

export const stripe = {
  createCardToken,
};
// enddoc: all
