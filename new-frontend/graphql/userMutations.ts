// graphql/userMutations.ts

import gql from 'graphql-tag';

export const UPDATE_USER_PROFILE = gql`
  mutation UpdateUserProfile($id: uuid!, $input: UserProfileInput!) {
    update_users_by_pk(
      pk_columns: { id: $id }
      _set: {
        name: $input.name
        bio: $input.bio
        location: $input.location
        website: $input.website
      }
    ) {
      id
      name
      bio
      location
      website
    }
  }
`;

export const UPDATE_USER_AVATAR = gql`
  mutation UpdateUserAvatar($id: uuid!, $avatar_url: String!) {
    update_users_by_pk(
      pk_columns: { id: $id }
      _set: { avatar_url: $avatar_url }
    ) {
      id
      avatar_url
    }
  }
`;