
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
// graphql/eventMutations.ts

export const FOLLOW_EVENT = gql`
  mutation FollowEvent($eventId: uuid!, $userId: uuid!) {
    insert_follows_one(object: { 
      event_id: $eventId,
      user_id: $userId
    }) {
      id
      created_at
    }
  }
`;


export const BOOKMARK_EVENT = gql`
  mutation BookmarkEvent($eventId: uuid!, $userId: uuid!) {
    insert_bookmarks_one(object: { 
      event_id: $eventId,
      user_id: $userId
    }) {
      id
      created_at
    }
  }
`;

export const UNBOOKMARK_EVENT = gql`
  mutation UnbookmarkEvent($eventId: uuid!) {
    delete_bookmarks(where: { event_id: { _eq: $eventId } }) {
      affected_rows
    }
  }
`;

export const UNFOLLOW_EVENT = gql`
  mutation UnfollowEvent($eventId: uuid!) {
    delete_follows(where: { event_id: { _eq: $eventId } }) {
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



export const UPDATE_EVENT_STATUS = gql`
  mutation UpdateEventStatus($eventId: uuid!, $status: String!) {
    update_events(
      where: { id: { _eq: $eventId } }
      _set: { status: $status }
    ) {
      affected_rows
      returning {
        id
        status
      }
    }
  }
`;


export const DELETE_EVENT = gql`
  mutation DeleteEvent($eventId: uuid!) {
     delete_events_by_pk(id: $eventId) {
      id
      title
    }
  }
`;


// graphql/eventMutations.ts

// Option 1: Using update_events_by_pk (most common)
export const UPDATE_EVENT_MUTATION = gql`
  mutation UpdateEvent($id: uuid!, $input: EventInput!) {
    update_events_by_pk(
      pk_columns: { id: $id }
      _set: {
        title: $input.title
        description: $input.description
        category_id: $input.category_id
        price: $input.price
        is_free: $input.is_free
        venue: $input.venue
        address: $input.address
        latitude: $input.latitude
        longitude: $input.longitude
        event_date: $input.event_date
        start_time: $input.start_time
        end_time: $input.end_time
        status: $input.status
        featured_image: $input.featured_image
        tags: $input.tags
      }
    ) {
      id
      title
      description
      featured_image
      is_free
      price
      event_date
      venue
      address
      latitude
      longitude
      status
      start_time
      end_time
      category_id
      tags
      updated_at
    }
  }
`;