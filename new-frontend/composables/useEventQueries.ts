// composables/useEventQueries.ts
import { useGraphQL } from '~/composables/useGraphQL'

export const useEventQueries = () => {
  const { query, mutate } = useGraphQL()

  const GET_EVENT_DETAILS = `
    query GetEventDetails($eventId: uuid!) {
      events_by_pk(id: $eventId) {
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
        featured_image
        status
        view_count
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
        event_images {
          id
          image_url
          is_featured
        }
        event_tags {
          tag {
            id
            name
          }
        }
      }
    }
  `

  const GET_EVENT_STATS = `
    query GetEventStats($eventId: uuid!) {
      get_event_stats(args: {event_id: $eventId}) {
        total_bookmarks
        total_follows
        total_tickets_sold
        total_revenue
      }
    }
  `

  const GET_USER_INTERACTIONS = `
    query GetUserInteractions($eventId: uuid!, $userId: uuid!) {
      bookmarks(
        where: {
          _and: [
            { event_id: { _eq: $eventId } },
            { user_id: { _eq: $userId } }
          ]
        }
      ) {
        id
      }
      follows(
        where: {
          _and: [
            { event_id: { _eq: $eventId } },
            { user_id: { _eq: $userId } }
          ]
        }
      ) {
        id
      }
    }
  `

  const TOGGLE_BOOKMARK = `
    mutation ToggleBookmark($eventId: uuid!, $userId: uuid!) {
      insert_bookmarks_one(
        object: { event_id: $eventId, user_id: $userId }
        on_conflict: {
          constraint: bookmarks_user_id_event_id_key
          update_columns: []
        }
      ) {
        id
      }
    }
  `

  const REMOVE_BOOKMARK = `
    mutation RemoveBookmark($eventId: uuid!, $userId: uuid!) {
      delete_bookmarks(
        where: {
          _and: [
            { event_id: { _eq: $eventId } },
            { user_id: { _eq: $userId } }
          ]
        }
      ) {
        affected_rows
      }
    }
  `

  const TOGGLE_FOLLOW = `
    mutation ToggleFollow($eventId: uuid!, $userId: uuid!) {
      insert_follows_one(
        object: { event_id: $eventId, user_id: $userId }
        on_conflict: {
          constraint: follows_user_id_event_id_key
          update_columns: []
        }
      ) {
        id
      }
    }
  `

  const REMOVE_FOLLOW = `
    mutation RemoveFollow($eventId: uuid!, $userId: uuid!) {
      delete_follows(
        where: {
          _and: [
            { event_id: { _eq: $eventId } },
            { user_id: { _eq: $userId } }
          ]
        }
      ) {
        affected_rows
      }
    }
  `

  const PURCHASE_TICKET = `
    mutation PurchaseTicket($eventId: uuid!, $userId: uuid!, $quantity: Int!, $totalPrice: numeric!) {
      insert_tickets_one(
        object: {
          event_id: $eventId
          user_id: $userId
          quantity: $quantity
          total_price: $totalPrice
          status: "confirmed"
        }
      ) {
        id
        quantity
        total_price
        status
        created_at
        event {
          title
          event_date
          venue
        }
      }
    }
  `

  const GET_USER_TICKETS = `
    query GetUserTickets($userId: uuid!) {
      tickets(
        where: { user_id: { _eq: $userId } }
        order_by: { created_at: desc }
      ) {
        id
        quantity
        total_price
        status
        created_at
        event {
          id
          title
          event_date
          venue
          featured_image
          is_free
        }
      }
    }
  `

  const GET_USER_BOOKMARKS = `
    query GetUserBookmarks($userId: uuid!) {
      bookmarks(
        where: { user_id: { _eq: $userId } }
        order_by: { created_at: desc }
      ) {
        id
        created_at
        event {
          id
          title
          event_date
          venue
          featured_image
          price
          is_free
        }
      }
    }
  `

  const GET_USER_FOLLOWS = `
    query GetUserFollows($userId: uuid!) {
      follows(
        where: { user_id: { _eq: $userId } }
        order_by: { created_at: desc }
      ) {
        id
        created_at
        event {
          id
          title
          event_date
          venue
          featured_image
          price
          is_free
        }
      }
    }
  `

  return {
    GET_EVENT_DETAILS,
    GET_EVENT_STATS,
    GET_USER_INTERACTIONS,
    TOGGLE_BOOKMARK,
    REMOVE_BOOKMARK,
    TOGGLE_FOLLOW,
    REMOVE_FOLLOW,
    PURCHASE_TICKET,
    GET_USER_TICKETS,
    GET_USER_BOOKMARKS,
    GET_USER_FOLLOWS
  }
}