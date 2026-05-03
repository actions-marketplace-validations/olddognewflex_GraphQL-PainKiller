import gql from "graphql-tag";

export const GET_DASHBOARD = gql`
  query GetUserDashboard {
    user {
      id
      assets {
        items {
          id
          vin
          buyer {
            accountId
            name
          }
          seller {
            accountId
            name
          }
          dealCharges {
            type
            amount
            status
          }
          inspections {
            status
            completedAt
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
