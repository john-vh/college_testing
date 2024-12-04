import { useCallback, useState } from 'react';
import { PostingInfo } from './useAllPostings.ts';

export function useActivatePosting() {
    const [loading, setLoading] = useState(false);

    const activatePosting = (post: PostingInfo) => {
        setLoading(true);
        const { business_id, id } = post;
        async function fetchData() {
            try {
                const response = await fetch(`http://127.0.0.1:8080/businesses/${business_id}/posts/${id}/activate`,
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

    return { activatePosting, loading };
}
