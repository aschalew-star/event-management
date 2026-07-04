// graphql/userQueries.ts

import gql from 'graphql-tag';

export const GET_USER_PROFILE = gql`
  query GetUserProfile($userId: uuid!) {
    users_by_pk(id: $userId) {
      id
      name
      email
      avatar_url
      bio
      location
      website
      created_at
      events_aggregate {
        aggregate {
          count
        }
      }
      bookmarks_aggregate {
        aggregate {
          count
        }
      }
      follows_aggregate {
        aggregate {
          count
        }
      }
      events(order_by: { created_at: desc }, limit: 3) {
        id
        title
        created_at
      }
      bookmarks(order_by: { created_at: desc }, limit: 2) {
        id
        created_at
        event {
          id
          title
        }
      }
    }
  }
`;