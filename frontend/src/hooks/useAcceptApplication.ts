import { PostingApplicationInfo } from './useApplicationInfo';

interface AcceptApplicationProps {
  entry: PostingApplicationInfo,
  index: number
}

export function useAcceptApplication() {

  const acceptApplication = (data: AcceptApplicationProps) => {
    const { business_id, post_id } = data.entry;
    const { id: user_id } = data.entry.applications[data.index].user;

    async function fetchData() {
      try {
        const response = await fetch(`${process.env.REACT_APP_API_URL}/businesses/${business_id}/posts/${post_id}/applications/${user_id}/accept`, { mode: "cors", credentials: 'include' });
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
      } catch (error) {
        console.log(error);
      }
    }
    fetchData();
  }; // Empty dependency array ensures this runs only once

  return acceptApplication;
}
