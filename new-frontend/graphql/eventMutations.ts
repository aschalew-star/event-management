import { gql } from 'graphql-tag'

export const CREATE_EVENT_MUTATION = gql`
  mutation CreateEventOne($input: CreateEventInput!) {
    createEvent(input: $input) {
      id
      message
      success
      featured_image
    }
  }
`