import { gql } from "graphql-tag";

export const CREATE_EVENT_MUTATION = gql`
  mutation CreateEventOne($input: CreateEventInput!) {
    createEvent(input: $input) {
      id
      message
      success
      featured_image
    }
  }
`;

export const RSVP_TO_EVENT = gql`
  mutation RSVPToEvent($eventId: uuid!) {
    insert_bookmarks_one(object: { event_id: $eventId }) {
      id
      event_id
      user_id
    }
  }
`;

export const UNRSVP_FROM_EVENT = gql`
  mutation UnRSVPFromEvent($eventId: uuid!) {
    delete_bookmarks(where: { event_id: { _eq: $eventId } }) {
      affected_rows
    }
  }
`;

export const ADD_EVENT_COMMENT = gql`
  mutation AddEventComment($eventId: uuid!, $content: String!) {
    insert_event_comments_one(
      object: { event_id: $eventId, content: $content }
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

export const GET_USER_BOOKMARKS = gql`
  query GetUserBookmarks($userId: uuid!) {
    bookmarks(where: { user_id: { _eq: $userId } }) {
      event_id
    }
  }
`;

export const UPDATE_EVENT_STATUS = gql`
  mutation UpdateEventStatus($eventId: uuid!, $status: String!) {
    update_events(where: { id: { _eq: $eventId } }, _set: { status: $status }) {
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

// If you need to handle tags separately, use this mutation instead
export const UPDATE_EVENT_WITH_TAGS = gql`
  mutation UpdateEvent(
    $id: uuid!
    $input: events_update_input!
    $tags: [String!]
  ) {
    update_events_by_pk(pk_columns: { id: $id }, _set: $input) {
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
      updated_at
    }
  }
`;

export const FOLLOW_EVENT = gql`
  mutation FollowEvent($followedUserId: uuid!) {
    insert_follows_one(object: { followed_user_id: $followedUserId }) {
      id
      follower_id
      followed_user_id
      created_at
    }
  }
`;

export const UNFOLLOW_EVENT = gql`
  mutation UnfollowEvent($followedUserId: uuid!) {
    delete_follows(where: { followed_user_id: { _eq: $followedUserId } }) {
      affected_rows
    }
  }
`;

export const BOOKMARK_EVENT = gql`
  mutation BookmarkEvent($eventId: uuid) {
    insert_bookmarks_one(object: { event_id: $eventId }) {
      id
      event_id
      user_id
      created_at
    }
  }
`;

export const UNBOOKMARK_EVENT = gql`
  mutation UnbookmarkEvent($eventId: uuid!, $userId: uuid!) {
    delete_bookmarks(
      where: { event_id: { _eq: $eventId }, user_id: { _eq: $userId } }
    ) {
      affected_rows
    }
  }
`;

export const CREATE_TICKET = gql`
  mutation CreateTicket(
    $eventId: uuid!
    $quantity: Int!
    $totalPrice: numeric!
    $status: String!
    $payment_id: uuid!
  ) {
    insert_tickets_one(
      object: {
        event_id: $eventId
        quantity: $quantity
        total_price: $totalPrice
        status: $status
        payment_id: $payment_id
      }
    ) {
      id
      event_id
      user_id
      quantity
      total_price
      status
      created_at
    }
  }
`;

export const CREATE_TICKET_WITH_PAYMENT = gql`
  mutation CreateTicketWithPayment(
    $eventId: uuid!
    $quantity: Int!
    $totalPrice: numeric!
    $status: String!
    $paymentId: String!
  ) {
    insert_tickets_one(
      object: {
        event_id: $eventId
        quantity: $quantity
        total_price: $totalPrice
        status: $status
        payment_id: $paymentId
      }
    ) {
      id
      event_id
      user_id
      quantity
      total_price
      status
      payment_id
      created_at
    }
  }
`;

export const CREATE_PAYMENT = gql`
  mutation CreatePayment {
    insert_payments_one(
      object: {
        status: "completed"
        amount: 0
        currency: "ETB"
        payment_method: "FREE"
      }
    ) {
      id
      user_id
      amount
      status
      payment_method
      created_at
    }
  }
`;

// new-frontend/graphql/eventMutations.ts

// eventMutations.ts

export const DELETE_EVENT_IMAGE_MUTATION = gql`
  mutation DeleteEventImage($imageId: uuid!) {
    delete_event_images_by_pk(id: $imageId) {
      id
      image_url
      event_id
    }
  }
`;

export const UPDATE_EVENT_MUTATION = gql`
  mutation UpdateEvent($id: uuid!, $input: events_set_input!) {
    update_events_by_pk(pk_columns: { id: $id }, _set: $input) {
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
    }
  }
`;

export const UPLOAD_EVENT_IMAGES_MUTATION = gql`
  mutation UploadEventImages(
    $eventId: uuid!
    $images: [event_images_insert_input!]!
  ) {
    insert_event_images(objects: $images) {
      returning {
        id
        image_url
        is_featured
        event_id
      }
    }
  }
`;

export const SET_FEATURED_IMAGE_MUTATION = gql`
  mutation SetFeaturedImage($eventId: uuid!, $imageId: uuid!) {
    # First, unset featured flag on all images for this event
    update_event_images(
      where: { event_id: { _eq: $eventId }, is_featured: { _eq: true } }
      _set: { is_featured: false }
    ) {
      affected_rows
    }

    # Then set the selected image as featured
    update_event_images_by_pk(
      pk_columns: { id: $imageId }
      _set: { is_featured: true }
    ) {
      id
      image_url
      is_featured
    }
  }
`;

export const UPLOAD_EVENT_IMAGES = gql`
  mutation UploadEventImages($eventId: uuid!, $images: [String!]!) {
    uploadEventImages(event_id: $eventId, images: $images) {
      success
      message
      urls
    }
  }
`;

export const DELETE_EVENT_IMAGE = gql`
  mutation delete_event_images($imageId: uuid!) {
    delete_event_images(where: { id: { _eq: $imageId } }) {
      affected_rows # Note: Hasura returns affected_rows instead of success/message by default
    }
  }
`;

export const SET_FEATURED_IMAGE = gql`
  mutation SetFeaturedImage($eventId: uuid!, $imageId: uuid!) {
    update_event_images(
      where: { event_id: { _eq: $eventId }, id: { _eq: $imageId } }
      _set: { is_featured: true }
    ) {
      affected_rows
      returning {
        id
        image_url
        is_featured
      }
    }
  }
`;


// if that's how your composable imports it!
export const removeFeaturedMutation = gql`
  mutation RemoveFeaturedImage($eventId: uuid!, $imageId: uuid!) {
    update_event_images(
      where: {
        event_id: { _eq: $eventId }
        id: { _eq: $imageId }
      }
      _set: {
        is_featured: false
      }
    ) {
      affected_rows
      returning {
        id
        image_url
        is_featured
      }
    }
  }
`;

export const UPDATE_EVENT = gql`
  mutation UpdateEvent($id: uuid!, $input: events_set_input!) {
    update_events_by_pk(pk_columns: { id: $id }, _set: $input) {
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
    }
  }
`;

// graphql/eventMutations.ts

export const FOLLOW_USER = gql`
  mutation FollowUser($followerId: uuid!, $followedUserId: uuid!) {
    insert_follows_one(object: {
      follower_id: $followerId
      followed_user_id: $followedUserId
    }) {
      id
      follower_id
      followed_user_id
      created_at
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
