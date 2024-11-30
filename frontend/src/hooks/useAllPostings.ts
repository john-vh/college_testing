import { useState, useEffect } from 'react';

export interface PostingInfo {
    id: number,
    business_id: string,
    created_at: string,
    updated_at: string,
    status: string,
    title: string,
    desc: string,
    pay: number,
    time_est: number
}

function useAllPostings(): PostingInfo[] | null {
    const [postingInfo, setPostingInfo] = useState<PostingInfo[] | null>(null);

    useEffect(() => {
        async function fetchData() {
            try {
                const response = await fetch('http://127.0.0.1:8080/posts', { mode: "cors", credentials: 'include' });
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                const data: PostingInfo[] = await response.json();
                setPostingInfo(data);
            } catch (error) {
                console.log(error);
            }
        }
        fetchData();
    }, []); // Empty dependency array ensures this runs only once

    return postingInfo;
}

export default useAllPostings;
