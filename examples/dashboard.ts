import gql from "graphql-tag";

export const GET_DASHBOARD = gql`
  query GetUserDashboard {
    user {
      id
      purchased {
        items {
          id
          vin
          seller {
            accountId
            name
          }
          charges {
            type
            amount
            status
          }
        }
      }
    }
  }
`;

export const GET_POSTS = /* GraphQL */ `
  query GetPostsFromCommentTemplate {
    posts {
      comments {
        author {
          name
        }
      }
    }
  }
`;
