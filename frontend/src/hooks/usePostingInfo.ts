import { useState, useCallback } from 'react';
import { PostingInfo } from './useAllPostings';
import { BusinessInfo } from './useBusinessInfo';

interface PostingInfoHook {
    data: PostingInfo[];
    business_map: Map<string, BusinessInfo>;
    loading: boolean;
    error: string | null;
    fetchPostingInfo: () => Promise<void>;
}

export function usePostingInfo(): PostingInfoHook {
    const [data, setData] = useState<PostingInfo[]>([]);
    const [business_map, setBusinessMap] = useState<Map<string, BusinessInfo>>(new Map());
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const fetchPostingInfo = useCallback(async () => {
        setLoading(true);
        setError(null);

        try {
            const response = await fetch('http://127.0.0.1:8080/users/0/posts',
                { mode: "cors", credentials: 'include' });
            const business_response = await fetch('http://127.0.0.1:8080/users/0/businesses',
                { mode: "cors", credentials: 'include' });
            if (!response.ok || !business_response.ok) {
                throw new Error('Network response was not ok');
            }

            const newData: PostingInfo[] = await response.json();
            const businessData: BusinessInfo[] = await business_response.json();
            const new_business_map = new Map<string, BusinessInfo>(
                businessData.map((obj) => [obj.id, obj])
            );
            setData(newData);
            setBusinessMap(new_business_map);
        }
        catch (error) {
            const errorMessage = error instanceof Error ? error.message : 'An error occurred';
            setError(errorMessage);
            setData([]);
            setBusinessMap(new Map<string, BusinessInfo>())
        } finally {
            setLoading(false);
        }
    }, []);

    // Initial fetch on mount
    useState(() => {
        fetchPostingInfo();
    });

    return {
        data,
        business_map,
        loading,
        error,
        fetchPostingInfo
    };
}