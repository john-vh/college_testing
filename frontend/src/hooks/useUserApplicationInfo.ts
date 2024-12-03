import { useState, useEffect, useMemo } from 'react';
import { PostingApplicationInfo } from './useApplicationInfo.ts';

export interface UserApplicationInfo {
    post: UserPostingInfo,
    business: UserBusinessInfo,
    status: string,
    created_at: string
}

interface UserPostingInfo {
    id: number,
    title: string,
    status: string,
    pay: number,
    time_est: number,
    updated_at: string,
    created_at: string
}

interface UserBusinessInfo {
    id: string,
    name: string,
    status: string,
    created_at: string
}

export function useUserApplicationInfo(): UserApplicationInfo[] {
    const [applicationInfo, setApplicationInfo] = useState<UserApplicationInfo[]>([]);

    useEffect(() => {
        async function fetchData() {
            try {
                const response = await fetch(`http://127.0.0.1:8080/users/0/applications`, { mode: "cors", credentials: 'include' });
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                setApplicationInfo(await response.json())
            } catch (error) {
                console.log(error);
            }
        }
        fetchData();
    }, []); // Empty dependency array ensures this runs only once

    return applicationInfo;
}