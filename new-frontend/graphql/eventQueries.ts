// graphql/eventQueries.ts
import { gql } from "@apollo/client/core";
import { graphql } from "graphql/graphql";

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

export const GET_ALL_EVENTS = gql`
  query GetAllEvents($filters: events_bool_exp!, $limit: Int!, $offset: Int!) {
    events(
      where: $filters
      limit: $limit
      offset: $offset
      order_by: [{ event_date: asc }]
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
      view_count
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

export const CHECK_EVENT_ATTENDANCE = gql`
  query CheckEventAttendance($eventId: uuid!, $userId: uuid!) {
    bookmarks(
      where: { event_id: { _eq: $eventId }, user_id: { _eq: $userId } }
    ) {
      id
    }
  }
`;

// // graphql/eventQueries.ts

// export const GET_USER_BOOKMARKS = gql`
//   query GetUserBookmarks($userId: uuid!) {
//     bookmarks(
//       where: { user_id: { _eq: $userId } }
//       order_by: { created_at: desc }
//     ) {
//       id
//       event_id
//       created_at
//       event {
//         id
//         title
//         featured_image
//         is_free
//         price
//         event_date
//         venue
//         user {
//           id
//           name
//         }
//       }
//     }
//   }
// `;

export const GET_USER_EVENTS_WITH_FOLLOWERS = gql`
  query GetUserEventsWithFollowers($userId: uuid!) {
    events(
      where: { user_id: { _eq: $userId } }
      order_by: { event_date: asc }
    ) {
      id
      title
      description
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
      updated_at
      user_id
      category_id

      # User relationship with follower count
      user {
        id
        name
        email
        avatar_url
        # Get followers count for the event creator
        followers_aggregate {
          aggregate {
            count
          }
        }
      }

      # Category relationship
      category {
        id
        name
        description
        icon
        color
      }

      # Event images
      event_images {
        id
        image_url
        is_featured
      }

      # Bookmarks aggregate
      bookmarks_aggregate {
        aggregate {
          count
        }
      }

      # Tickets aggregate (to show ticket sales)
      tickets_aggregate(where: { status: { _eq: "confirmed" } }) {
        aggregate {
          count
        }
      }
    }
  }
`;

// GraphQL Mutations
export const FOLLOW_EVENT = gql`
  mutation FollowEvent($eventId: uuid!, $userId: uuid!) {
    insert_follows_one(
      object: { followed_user_id: $eventId, follower_id: $userId }
    ) {
      id
    }
  }
`;

export const UNFOLLOW_EVENT = gql`
  mutation UnfollowEvent($eventId: uuid!) {
    delete_follows(where: { followed_user_id: { _eq: $eventId } }) {
      affected_rows
    }
  }
`;

export const BOOKMARK_EVENT = gql`
  mutation BookmarkEvent($eventId: uuid!, $userId: uuid!) {
    insert_bookmarks_one(object: { event_id: $eventId, user_id: $userId }) {
      id
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

export const CREATE_TICKET = gql`
  mutation CreateTicket(
    $eventId: uuid!
    $userId: uuid!
    $quantity: Int!
    $totalPrice: numeric!
    $status: String!
  ) {
    insert_tickets_one(
      object: {
        event_id: $eventId
        user_id: $userId
        quantity: $quantity
        total_price: $totalPrice
        status: $status
      }
    ) {
      id
      status
    }
  }
`;

export const GET_RELATED_EVENTS = gql`
  query GetRelatedEvents($categoryId: uuid!, $eventId: uuid!, $limit: Int!) {
    events(
      where: {
        category_id: { _eq: $categoryId }
        id: { _neq: $eventId }
        status: { _eq: "published" }
      }
      limit: $limit
      order_by: { created_at: desc }
    ) {
      id
      title
      venue
      event_date
    }
  }
`;

export const GET_EVENT_ATTENDEES = gql`
  query GetEventAttendees($eventId: uuid!) {
    tickets(
      where: { event_id: { _eq: $eventId }, status: { _eq: "confirmed" } }
      distinct_on: user_id
    ) {
      id
      user {
        id
        name
        avatar_url
      }
    }
  }
`;

export const CHECK_EVENT_FOLLOW = gql`
  query CheckEventFollow($userId: uuid!, $followedUserId: uuid!) {
    follows(
      where: {
        follower_id: { _eq: $userId }
        followed_user_id: { _eq: $followedUserId }
      }
    ) {
      id
      follower_id
      followed_user_id
    }
  }
`;

export const CHECK_EVENT_BOOKMARK = gql`
  query CheckEventBookmark($eventId: uuid!, $userId: uuid!) {
    bookmarks(
      where: { event_id: { _eq: $eventId }, user_id: { _eq: $userId } }
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
    follows_aggregate(where: { followed_user_id: { _eq: $eventId } }) {
      aggregate {
        count
      }
    }
  }
`;

export const GET_USER_TICKET = gql`
  query GetUserTicket($userId: uuid!, $eventId: uuid!) {
    tickets(
      where: {
        user_id: { _eq: $userId }
        event_id: { _eq: $eventId }
        status: { _eq: "confirmed" }
      }
    ) {
      id
      status
      event_id
    }
  }
`;

export const GET_EVENT_IMAGES = gql`
  query GetEventImages($eventId: uuid!) {
    event_images(where: { event_id: { _eq: $eventId } }) {
      id
      image_url
      is_featured
      event_id
    }
  }
`;

export const GET_EVENT_IMAGES_BY_EVENT_IDS = gql`
  query GetEventImagesByEventIds($eventIds: [uuid!]!) {
    event_images(
      where: { event_id: { _in: $eventIds }, is_featured: { _eq: true } }
      distinct_on: event_id
    ) {
      event_id
      image_url
    }
  }
`;

export const GET_EVENT_BY_ID = gql`
  query GetEventById($id: uuid!) {
    events_by_pk(id: $id) {
      id
      title
      description
      category_id
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
      event_images(order_by: { created_at: asc }) {
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
      tickets_aggregate(where: { status: { _eq: "confirmed" } }) {
        aggregate {
          count
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
      created_at
      user_id
      event_id
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

        # User relationship
        user {
          id
          name
          email
          avatar_url
        }

        # Category relationship
        category {
          id
          name
          icon
          color
        }

        # Event images
        event_images {
          id
          image_url
          is_featured
        }

        # Bookmarks aggregate
        bookmarks_aggregate {
          aggregate {
            count
          }
        }
      }
    }
  }
`;

// graphql/eventQueries.ts

export const GET_POPULAR_EVENTS = gql`
  query GetPopularEvents {
    events(
      where: { status: { _eq: "published" } }
      order_by: { view_count: desc }
      limit: 6
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
      view_count
      created_at
      updated_at
      user_id
      category_id

      user {
        id
        name
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

      tickets_aggregate(where: { status: { _eq: "confirmed" } }) {
        aggregate {
          count
        }
      }
    }
  }
`;


export const GET_USER_PROFILE = gql`
  query GET_USER_PROFILE($id: uuid!) {
    users_by_pk(id: $id) {
      id
      name
      bio
      avatar_url
      created_at
      # 1. FIXED: Changed from followers_aggregate to follows_aggregate
      follows_aggregate {
        aggregate {
          count
        }
      }
      # 2. FIXED: Changed from following_aggregate to followsByFollowerId_aggregate
      followsByFollowerId_aggregate {
        aggregate {
          count
        }
      }
    }
  }
`;

export const GET_USER_EVENTS = gql`
  query GetUserEvents($userId: uuid!) {
    events(
      where: { user_id: { _eq: $userId }, status: { _neq: "draft" } }
      order_by: { created_at: desc }
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
      view_count
      created_at
      updated_at
      category_id
      user_id

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

      tickets_aggregate(where: { status: { _eq: "confirmed" } }) {
        aggregate {
          count
        }
      }
    }
  }
`;

export const GET_USER_FOLLOWERS = gql`
  query GetUserFollowers($userId: uuid!) {
    follows(
      where: { followed_user_id: { _eq: $userId } }
      order_by: { created_at: desc }
    ) {
      id
      follower_id
      followed_user_id
      created_at
      userByFollowerId {
        id
        name
        email
        avatar_url
        bio
        created_at
      }
    }
  }
`;

// graphql/eventQueries.ts

export const GET_USER_PROFILES = gql`
  query GetUserProfile($id: uuid!) {
    users_by_pk(id: $id) {
      id
      name
      email
      avatar_url
      bio
      created_at
      followers_aggregate {
        aggregate {
          count
        }
      }
      following_aggregate {
        aggregate {
          count
        }
      }
    }
  }
`;

export const GET_USER_EVENTSS = gql`
  query GetUserEvents($userId: uuid!) {
    events(
      where: { user_id: { _eq: $userId }, status: { _neq: "draft" } }
      order_by: { created_at: desc }
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
      view_count
      created_at
      updated_at
      category_id
      user_id

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

      tickets_aggregate(where: { status: { _eq: "confirmed" } }) {
        aggregate {
          count
        }
      }
    }
  }
`;

export const GET_FOLLOW_STATUS = gql`
  query GetFollowStatus($followerId: uuid!, $followedUserId: uuid!) {
    follows(
      where: {
        follower_id: { _eq: $followerId }
        followed_user_id: { _eq: $followedUserId }
      }
    ) {
      id
      follower_id
      followed_user_id
      created_at
    }
  }
`;
// graphql/eventQueries.ts

export const GET_USER_FOLLOWING = gql`
  query GetUserFollowing($userId: uuid!) {
    follows(
      where: { follower_id: { _eq: $userId } }
      order_by: { created_at: desc }
    ) {
      id
      follower_id
      followed_user_id
      created_at
      userByFollowedUserId {
        id
        name
        email
        avatar_url
        bio
        created_at
        events(where: { status: { _neq: "draft" } }) {
          id
          title
        }
        events_aggregate(where: { status: { _neq: "draft" } }) {
          aggregate {
            count
          }
        }
        followers_aggregate {
          aggregate {
            count
          }
        }
      }
    }
  }
`;

export const GET_USER_FOLLOWINGs = gql`
  query GetUserFollowing($userId: uuid!) {
    follows(
      where: { follower_id: { _eq: $userId } }
      order_by: { created_at: desc }
    ) {
      id
      follower_id
      followed_user_id
      created_at
      user {
        id
        name
        email
        avatar_url
        bio
        created_at
        events(where: { status: { _neq: "draft" } }) {
          id
          title
        }
      }
    }
  }
`;

export const UNFOLLOW_USER = gql`
  mutation UnfollowUser($followerId: uuid!, $followedUserId: uuid!) {
    delete_follows(
      where: {
        follower_id: { _eq: $followerId }
        followed_user_id: { _eq: $followedUserId }
      }
    ) {
      affected_rows
    }
  }
`;
