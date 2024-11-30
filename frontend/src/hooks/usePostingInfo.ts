import { useState, useCallback } from 'react';
import { PostingInfo } from './useAllPostings';

interface PostingInfoHook {
    data: PostingInfo[];
    loading: boolean;
    error: string | null;
    fetchPostingInfo: () => Promise<void>;
}

export function usePostingInfo(): PostingInfoHook {
    const [data, setData] = useState<PostingInfo[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const fetchPostingInfo = useCallback(async () => {
        setLoading(true);
        setError(null);

        try {
            const response = await fetch('http://127.0.0.1:8080/users/0/posts',
                { mode: "cors", credentials: 'include' });

            if (!response.ok) {
                throw new Error('Network response was not ok');
            }

            const newData = await response.json();
            setData(newData);
        }
        catch (error) {
            const errorMessage = error instanceof Error ? error.message : 'An error occurred';
            setError(errorMessage);
            setData([]);
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
        loading,
        error,
        fetchPostingInfo
    };
}