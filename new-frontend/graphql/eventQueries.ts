// graphql/eventQueries.ts
import { gql } from '@apollo/client/core';

export const GET_ALL_EVENTS = gql`
  query GetAllEvents($filters: events_bool_exp, $limit: Int, $offset: Int) {
    events(
      where: $filters
      limit: $limit
      offset: $offset
      order_by: { event_date: asc }
    ) {
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
      featured_image
      created_at
      category {
        id
        name
        icon
        color
      }
      user {
        id
        name
        email
        avatar_url
      }
    }
    events_aggregate(where: $filters) {
      aggregate {
        count
      }
    }
  }
`;

export const GET_EVENT_CATEGORIES = gql`
  query GetEventCategories {
    categories {
      id
      name
      icon
      color
      description
    }
  }
`;

export const GET_UPCOMING_EVENTS = gql`
  query GetUpcomingEvents($limit: Int) {
    events(
      where: { event_date: { _gte: "now()" }, status: { _eq: "published" } }
      limit: $limit
      order_by: { event_date: asc }
    ) {
      id
      title
      venue
      event_date
      featured_image
      is_free
      price
      category {
        name
        color
      }
    }
  }
`;

export const GET_FEATURED_EVENTS = gql`
  query GetFeaturedEvents($limit: Int) {
    events(
      where: { status: { _eq: "published" } }
      limit: $limit
      order_by: { created_at: desc }
    ) {
      id
      title
      description
      venue
      event_date
      featured_image
      is_free
      price
      category {
        name
        color
        icon
      }
    }
  }
`;

export const GET_EVENT_COMMENTS = gql`
  query GetEventComments($eventId: uuid!, $limit: Int = 20) {
    comments(
      where: { event_id: { _eq: $eventId } }
      order_by: { created_at: desc }
      limit: $limit
    ) {
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


export const GET_EVENT_BY_ID = gql`
  query GetEventById($id: uuid!) {
    events_by_pk(id: $id) {
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
      featured_image
      view_count
      created_at
      updated_at
      category {
        id
        name
        icon
        color
      }
      user {
        id
        name
        email
        avatar_url
      }
    }
  }
`;

export const GET_RELATED_EVENTS = gql`
  query GetRelatedEvents($categoryId: uuid, $eventId: uuid, $limit: Int = 4) {
    events(
      where: {
        category_id: { _eq: $categoryId }
        id: { _neq: $eventId }
        status: { _eq: "published" }
      }
      limit: $limit
      order_by: { event_date: asc }
    ) {
      id
      title
      description
      venue
      event_date
      featured_image
      is_free
      price
      category {
        name
        color
        icon
      }
    }
  }
`;

export const GET_EVENT_ATTENDEES = gql`
  query GetEventAttendees($eventId: uuid!) {
    bookmarks(where: { event_id: { _eq: $eventId } }) {
      id
      user {
        id
        name
        avatar_url
      }
    }
  }
`;

export const CHECK_EVENT_ATTENDANCE = gql`
  query CheckEventAttendance($eventId: uuid!, $userId: uuid!) {
    bookmarks(
      where: {
        event_id: { _eq: $eventId }
        user_id: { _eq: $userId }
      }
    ) {
      id
    }
  }
`;

export const CHECK_EVENT_FOLLOW = gql`
  query CheckEventFollow($eventId: uuid!, $userId: uuid!) {
    follows(
      where: {
        event_id: { _eq: $eventId }
        user_id: { _eq: $userId }
      }
    ) {
      id
    }
  }
`;

export const CHECK_EVENT_BOOKMARK = gql`
  query CheckEventBookmark($eventId: uuid!, $userId: uuid!) {
    bookmarks(
      where: {
        event_id: { _eq: $eventId }
        user_id: { _eq: $userId }
      }
    ) {
      id
    }
  }
`;

export const GET_EVENT_STATS = gql`
  query GetEventStats($eventId: uuid!) {
    bookmarks_aggregate(where: { event_id: { _eq: $eventId } }) {
      aggregate {
        count
      }
    }
    follows_aggregate(where: { event_id: { _eq: $eventId } }) {
      aggregate {
        count
      }
    }
  }
`;

export const GET_USER_TICKET = gql`
  query GetUserTickets($userId: uuid!) {
    tickets(
      where: {
        user_id: { _eq: $userId }
        status: { _eq: "confirmed" }
      }
    ) {
      id
      event_id
      status
    }
  }
`;




// graphql/eventQueries.ts

export const GET_USER_FOLLOWS = gql`
  query GetUserFollows($userId: uuid!) {
    follows(
      where: { user_id: { _eq: $userId } }
      order_by: { created_at: desc }
    ) {
      id
      created_at
      event_id
      user_id
      event {
        id
        title
        description
        featured_image
        is_free
        price
        event_date
        venue
        address
        status
        user {
          id
          name
          email
        }
      }
    }
  }
`;

// graphql/eventQueries.ts

export const GET_USER_BOOKMARKS = gql`
  query GetUserBookmarks($userId: uuid!) {
    bookmarks(
      where: { user_id: { _eq: $userId } }
      order_by: { created_at: desc }
    ) {
      id
      event_id
      created_at
      event {
        id
        title
        featured_image
        is_free
        price
        event_date
        venue
        user {
          id
          name
        }
      }
    }
  }
`;

// graphql/eventQueries.ts


// graphql/eventQueries.ts

export const GET_USER_EVENTS = gql`
  query GetUserEvents($userId: uuid!) {
    events(
      where: { 
        user_id: { _eq: $userId }
        status: { _neq: "draft" }
      }
      order_by: { event_date: asc }
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
      view_count
      created_at
      user_id
      category_id
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
      }
    }
  }
`;

// graphql/eventQueries.ts

export const GET_USER_FOLLOWERS = gql`
  query GetUserFollowers($userId: uuid!) {
    follows(
      where: { user_id: { _eq: $userId } }
      order_by: { created_at: desc }
    ) {
      id
      created_at
      user_id
      event_id
      user {
        id
        name
        email
        avatar_url
        bio
      }
      event {
        id
        title
        featured_image
        event_date
        venue
        is_free
        price
        status
      }
    }
  }
`;
