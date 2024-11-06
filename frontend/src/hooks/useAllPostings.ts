import { useState, useEffect } from 'react';

interface PostingInfo {
    id: number,
    business_id: number,
    created_at: string,
    updated_at: string,
    status: number,
    title: string,
    desc: string
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
