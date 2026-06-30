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
        avatar
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