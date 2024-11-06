import { useState, useEffect } from 'react';

interface PostingInfo {
    id: string,
    status: number,
    title: string,
    desc: string,
    business_id: string,
    created_at: string,
    updated_at?: string,
}

interface PostingInfoProps {
    business_ids: string[]
}

export function usePostingInfo({ business_ids }: PostingInfoProps): PostingInfo[] {
    const [postingInfo, setPostingInfo] = useState<PostingInfo[]>([]);

    useEffect(() => {
        async function fetchData() {
            for (const business_id in business_ids) {
                try {
                    const response = await fetch(`http://127.0.0.1:8080/businesses/${business_id}/posts`, { mode: "cors", credentials: 'include' });
                    if (!response.ok) {
                        throw new Error('Network response was not ok');
                    }
                    setPostingInfo([...postingInfo, await response.json()]);
                } catch (error) {
                    console.log(error);
                }
            }
        }
        fetchData();
    }, [business_ids, postingInfo]); // Empty dependency array ensures this runs only once

    return postingInfo;
}