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

interface PostingInfoProps {
  isAdmin: boolean;
}

export function usePostingInfo({ isAdmin }: PostingInfoProps): PostingInfoHook {
  const [data, setData] = useState<PostingInfo[]>([]);
  const [business_map, setBusinessMap] = useState<Map<string, BusinessInfo>>(new Map());
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchPostingInfo = useCallback(async () => {
    setLoading(true);
    setError(null);

    try {
      let response;
      let business_response;
      if (isAdmin) {
        response = await fetch(`${process.env.REACT_APP_API_URL}/admin/posts`,
          { mode: "cors", credentials: 'include' });
        business_response = await fetch(`${process.env.REACT_APP_API_URL}/admin/businesses`,
          { mode: "cors", credentials: 'include' });
      }
      else {
        response = await fetch(`${process.env.REACT_APP_API_URL}/users/0/posts`,
          { mode: "cors", credentials: 'include' });
        business_response = await fetch(`${process.env.REACT_APP_API_URL}/users/0/businesses`,
          { mode: "cors", credentials: 'include' });
      }

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
  }, [isAdmin]);

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
