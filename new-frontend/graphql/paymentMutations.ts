
import { gql } from '@apollo/client/core'

export const PROCESS_CHAPA_PAYMENT = gql`
  mutation ProcessChapaPayment($input: ProcessChapaPaymentInput!) {
    processChapaPayment(input: $input) {
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
  mutation VerifyChapaPayment($input: VerifyChapaPaymentInput!) {
    verifyChapaPayment(input: $input) {
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