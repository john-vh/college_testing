import { useState, useEffect } from 'react';

interface ApplicationInfo {
    id: string,
    status: number,
    title: string,
    desc: string,
    business_id: string,
    created_at: string,
    updated_at?: string,
}

interface ApplicationInfoProps {
    business_ids: string[]
}

export function useApplicationInfo({ business_ids }: ApplicationInfoProps): ApplicationInfo[] {
    const [applicationInfo, setApplicationInfo] = useState<ApplicationInfo[]>([]);

    useEffect(() => {
        async function fetchData() {
            for (const business_id in business_ids) {
                try {
                    const response = await fetch(`http://127.0.0.1:8080/businesses/${business_id}/posts/1/applications`, { mode: "cors", credentials: 'include' });
                    if (!response.ok) {
                        throw new Error('Network response was not ok');
                    }
                    setApplicationInfo([...applicationInfo, await response.json()]);
                } catch (error) {
                    console.log(error);
                }
            }
        }
        fetchData();
    }, [applicationInfo, business_ids]); // Empty dependency array ensures this runs only once

    return applicationInfo;
}