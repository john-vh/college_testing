import { useState, useEffect } from 'react';
import { usePostingIds } from './usePostingIds';
import { AccountInfo } from './useAccountInfo';

export interface ApplicationInfo {
  user: AccountInfo,
  notes: string,
  status: number
}

export interface PostingApplicationInfo {
  business_id: string,
  post_id: number,
  applications: ApplicationInfo[]
}

export function useApplicationInfo(): PostingApplicationInfo[] {
  const [applicationInfo, setApplicationInfo] = useState<PostingApplicationInfo[]>([]);
  const post_ids = usePostingIds();

  useEffect(() => {
    async function fetchData() {
      const allData: PostingApplicationInfo[] = [];
      for (const [business_id, post_id] of post_ids) {
        try {
          const response = await fetch(`${process.env.REACT_APP_API_URL}/businesses/${business_id}/posts/${post_id}/applications`, { mode: "cors", credentials: 'include' });
          if (!response.ok) {
            throw new Error('Network response was not ok');
          }
          allData.push(await response.json());
        } catch (error) {
          console.log(error);
        }
      }
      setApplicationInfo(allData);
    }
    fetchData();
  }, [post_ids]); // Empty dependency array ensures this runs only once

  return applicationInfo;
}
