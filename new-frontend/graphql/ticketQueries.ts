// graphql/ticketQueries.ts
import { gql } from '@apollo/client/core';

export const GET_USER_TICKETS = gql`
  query GetUserTickets($userId: uuid!) {
    tickets(
      where: { user_id: { _eq: $userId } }
      order_by: { created_at: desc }
    ) {
      id
      event_id
      user_id
      quantity
      total_price
      status
      payment_id
      transaction_ref
      created_at
      updated_at
      event {
        id
        title
        featured_image
        venue
        address
        event_date
        price
        is_free
        user {
          id
          name
        }
      }
    }
  }
`;

export const GET_TICKET_BY_ID = gql`
  query GetTicketById($id: uuid!) {
    tickets_by_pk(id: $id) {
      id
      event_id
      user_id
      quantity
      total_price
      status
      payment_id
      transaction_ref
      created_at
      updated_at
      event {
        id
        title
        description
        featured_image
        venue
        address
        event_date
        start_time
        end_time
        price
        is_free
        user {
          id
          name
          email
        }
      }
      user {
        id
        name
        email
      }
    }
  }
`;

export const GET_EVENT_TICKET_STATS = gql`
  query GetEventTicketStats($eventId: uuid!) {
    tickets_aggregate(
      where: { 
        event_id: { _eq: $eventId }
        status: { _eq: "confirmed" }
      }
    ) {
      aggregate {
        count
        sum {
          total_price
        }
      }
    }
  }
`;