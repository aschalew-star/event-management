
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

export const RSVP_TO_EVENT = gql`
  mutation RSVPToEvent($eventId: uuid!) {
    insert_bookmarks_one(object: {
      event_id: $eventId
    }) {
      id
      event_id
      user_id
    }
  }
`;

export const UNRSVP_FROM_EVENT = gql`
  mutation UnRSVPFromEvent($eventId: uuid!) {
    delete_bookmarks(
      where: {
        event_id: { _eq: $eventId }
      }
    ) {
      affected_rows
    }
  }
`;

export const ADD_EVENT_COMMENT = gql`
  mutation AddEventComment($eventId: uuid!, $content: String!) {
    insert_event_comments_one(object: {
      event_id: $eventId
      content: $content
    }) {
      id
      content
      created_at
      user {
        id
        name
        avatar_url
      }
    }
  }
`;

export const GET_USER_BOOKMARKS = gql`
  query GetUserBookmarks($userId: uuid!) {
    bookmarks(where: { user_id: { _eq: $userId } }) {
      event_id
    }
  }
`;


export const FOLLOW_EVENT = gql`
  mutation FollowEvent($eventId: uuid!) {
    insert_follows_one(object: {
      event_id: $eventId
    }) {
      id
    }
  }
`;

export const UNFOLLOW_EVENT = gql`
  mutation UnfollowEvent($eventId: uuid!) {
    delete_follows(
      where: {
        event_id: { _eq: $eventId }
      }
    ) {
      affected_rows
    }
  }
`;

export const BOOKMARK_EVENT = gql`
  mutation BookmarkEvent($eventId: uuid!) {
    insert_bookmarks_one(object: {
      event_id: $eventId
    }) {
      id
    }
  }
`;

export const UNBOOKMARK_EVENT = gql`
  mutation UnbookmarkEvent($eventId: uuid!) {
    delete_bookmarks(
      where: {
        event_id: { _eq: $eventId }
      }
    ) {
      affected_rows
    }
  }
`;

export const CREATE_TICKET = gql`
  mutation CreateTicket($eventId: uuid!, $quantity: Int!, $totalPrice: numeric!, $status: String!) {
    insert_tickets_one(object: {
      event_id: $eventId
      quantity: $quantity
      total_price: $totalPrice
      status: $status
    }) {
      id
      status
      created_at
    }
  }
`;