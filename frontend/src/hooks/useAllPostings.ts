import { useState, useEffect } from 'react';
import { BusinessInfo } from './useBusinessInfo';

export interface PostingInfo {
    id: number,
    business_id: string,
    created_at: string,
    updated_at: string,
    status: string,
    title: string,
    desc: string,
    pay: number,
    time_est: number,
    business_name?: string,
    business_website?: string,
    business_desc?: string
}

function useAllPostings() {
    const [postingInfo, setPostingInfo] = useState<PostingInfo[]>([]);
    const [businessMap, setBusinessMap] = useState<Map<string, BusinessInfo>>(new Map());

    useEffect(() => {
        async function fetchData() {
            try {
                const business_response = await fetch('http://127.0.0.1:8080/businesses', { mode: "cors", credentials: 'include' });
                const response = await fetch('http://127.0.0.1:8080/posts', { mode: "cors", credentials: 'include' });
                if (!response.ok || !business_response.ok) {
                    throw new Error('Network response was not ok');
                }
                const data: PostingInfo[] = await response.json();
                const businessData: BusinessInfo[] = await business_response.json();
                const new_business_map = new Map<string, BusinessInfo>(businessData.map((obj) => [obj.id, obj]));
                setPostingInfo(data);
                setBusinessMap(new_business_map);
            } catch (error) {
                console.log(error);
            }
        }
        fetchData();
    }, []); // Empty dependency array ensures this runs only once

    return { postingInfo, businessMap };
}

export default useAllPostings;
