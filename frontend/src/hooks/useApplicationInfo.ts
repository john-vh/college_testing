import { useState, useEffect } from 'react';
import { usePostingIds } from './usePostingIds.ts';
import { AccountInfo } from './useAccountInfo.ts';

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

interface ApplicationInfoProps {
  isAdmin: boolean;
}

export function useApplicationInfo({ isAdmin }: ApplicationInfoProps): PostingApplicationInfo[] {
  const [applicationInfo, setApplicationInfo] = useState<PostingApplicationInfo[]>([]);
  const post_ids = usePostingIds(isAdmin);

  useEffect(() => {
    async function fetchData() {
      const allData: PostingApplicationInfo[] = [];
      for (const [business_id, post_id] of post_ids) {
        try {
          let response;
          if (isAdmin) {
            response = await fetch(`${process.env.REACT_APP_API_URL}/businesses/${business_id}/posts/${post_id}/applications`, { mode: "cors", credentials: 'include' });
          }
          response = await fetch(`${process.env.REACT_APP_API_URL}/businesses/${business_id}/posts/${post_id}/applications`, { mode: "cors", credentials: 'include' });

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
