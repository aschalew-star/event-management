import { gql } from 'graphql-tag'

export const PROCESS_CHAPA_PAYMENT = gql`
  mutation ProcessChapaPayment(
    $eventId: String!
    $quantity: Int!
    $email: String!
    $phone: String!
    $firstName: String!
    $lastName: String!
  ) {
    processChapaPayment(
      input: {
        eventId: $eventId
        quantity: $quantity
        email: $email
        phone: $phone
        firstName: $firstName
        lastName: $lastName
      }
    ) {
      success
      message
      data {
        checkout_url
        transaction_ref
      }
    }
  }
`

export const VERIFY_CHAPA_PAYMENT = gql`
  mutation VerifyChapaPayment($transactionRef: String!) {
    verifyChapaPayment(input: { transactionRef: $transactionRef }) {
      success
      message
      status
      ticket {
        id
        event_id
        quantity
        total_price
        status
      }
    }
  }
`