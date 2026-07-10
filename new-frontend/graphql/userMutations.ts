import gql from 'graphql-tag';

export const UPDATE_USER_PROFILE = gql`
  mutation UpdateUserProfile($id: uuid!, $input: users_set_input!) {
    update_users_by_pk(
      pk_columns: { id: $id }
      _set: $input
    ) {
      id
      name
      bio
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