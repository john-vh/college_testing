import { useState } from 'react';
import { PostingInfo } from './useAllPostings';

export function useDeactivatePosting() {
  const [loading, setLoading] = useState(false);

  const deactivatePosting = (post: PostingInfo) => {
    setLoading(true);
    const { business_id, id } = post;
    async function fetchData() {
      try {
        const response = await fetch(`${process.env.REACT_APP_API_URL}/businesses/${business_id}/posts/${id}/deactivate`,
          { method: "POST", mode: "cors", credentials: 'include' });
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
      } catch (error) {
        console.log(error);
      }
      finally {
        setLoading(false);
      }
    }
    fetchData();
  };

  return { deactivatePosting, loading };
}
