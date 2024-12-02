import { useState, useEffect } from 'react';

interface PostingInfo {
    id: number,
    status: number,
    title: string,
    desc: string,
    business_id: string,
    created_at: string,
    updated_at?: string,
}

export function usePostingInfo(): PostingInfo[] {
    const [postingInfo, setPostingInfo] = useState<PostingInfo[]>([]);

    useEffect(() => {
        async function fetchData() {
            try {
                const response = await fetch(`http://127.0.0.1:8080/users/0/posts`, { mode: "cors", credentials: 'include' });
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                setPostingInfo(await response.json());
            } catch (error) {
                console.log(error);
            }

        }
        fetchData();
    }, []); // Empty dependency array ensures this runs only once

    return postingInfo;
}