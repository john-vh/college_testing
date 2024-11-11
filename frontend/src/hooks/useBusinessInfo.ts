import { useState, useEffect } from 'react';

interface BusinessInfo {
    id: string,
    created_at: string,
    status: number,
    name: string,
    desc: string,
    website: string
}

export function useBusinessInfo(): BusinessInfo[] {
    const [businessInfo, setBusinessInfo] = useState<BusinessInfo[]>([]);

    useEffect(() => {
        async function fetchData() {
            try {
                const response = await fetch("http://127.0.0.1:8080/users/0/businesses", { mode: "cors", credentials: 'include' });
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                const data: BusinessInfo[] = await response.json();
                setBusinessInfo(data);
            } catch (error) {
                console.log(error);
            }
        }
        fetchData();
    }, []); // Empty dependency array ensures this runs only once

    return businessInfo;
}