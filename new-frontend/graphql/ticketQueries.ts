// graphql/ticketQueries.ts
import { gql } from '@apollo/client/core';
  

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


export const GET_USER_TICKETS = gql`
  query GetUserTickets($userId: uuid!) {
    tickets(
      where: { user_id: { _eq: $userId } }
      order_by: { created_at: desc }
    ) {
      id
      event_id
      user_id
      payment_id
      quantity
      total_price
      status
      created_at
      updated_at
      event {
        id
        title
        description
        price
        is_free
        venue
        address
        latitude
        longitude
        event_date
        start_time
        end_time
        status
        view_count
        created_at
        updated_at
        category_id
        user_id
        user {
          id
          name
          email
          avatar_url
        }
        category {
          id
          name
          icon
          color
        }
        event_images {
          id
          image_url
          is_featured
          created_at
        }
        bookmarks_aggregate {
          aggregate {
            count
          }
        }
      }
      payment {
        id
        amount
        currency
        status
        transaction_ref
        payment_method
        created_at
      }
    }
  }
`;