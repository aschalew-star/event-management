// graphql/events.ts
import gql from 'graphql-tag'

// Hasura mutation for event creation (if using GraphQL directly)
export const INSERT_EVENT_MUTATION = gql`
  mutation InsertEvent($object: events_insert_input!) {
    insert_events_one(object: $object) {
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
      created_at
      updated_at
    }
  }
`

// Get categories and metadata
export const GET_EVENT_METADATA_QUERY = gql`
  query GetEventMetadata {
    categories {
      id
      name
      description
    }
}
`

// Get single event with all data
export const GET_EVENT_QUERY = gql`
  query GetEvent($id: uuid!) {
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
      images
      created_at
      updated_at
      category {
        id
        name
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

// Get all events with filters
export const GET_EVENTS_QUERY = gql`
  query GetEvents($where: events_bool_exp, $orderBy: [events_order_by!]) {
    events(where: $where, order_by: $orderBy) {
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
      images
      created_at
      updated_at
      category {
        id
        name
      }
    }
  }
`