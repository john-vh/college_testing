import { useState, useEffect, useCallback } from 'react';

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

export function useUserApplicationInfo() {
  const [applicationInfo, setApplicationInfo] = useState<UserApplicationInfo[]>([]);

  const fetchData = useCallback(async () => {
    try {
      const response = await fetch(`${process.env.REACT_APP_API_URL}/users/0/applications`, { mode: "cors", credentials: 'include' });
      if (!response.ok) {
        throw new Error('Network response was not ok');
      }
      setApplicationInfo(await response.json())
    } catch (error) {
      console.log(error);
    }
  }, []); // Empty dependency array ensures this runs only once

  useState(() => {
    fetchData();
  });

  return { applicationInfo, fetchData };
}
