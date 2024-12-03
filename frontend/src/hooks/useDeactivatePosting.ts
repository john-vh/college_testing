import { useState } from 'react';
import { PostingInfo } from './useAllPostings.ts';

export function useDeactivatePosting() {
    const [loading, setLoading] = useState(false);

    const deactivatePosting = (post: PostingInfo) => {
        setLoading(true);
        const { business_id, id } = post;
        async function fetchData() {
            try {
                const response = await fetch(`http://127.0.0.1:8080/businesses/${business_id}/posts/${id}/deactivate`,
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
