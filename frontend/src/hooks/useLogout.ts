import { useCallback } from 'react';

export function useLogout() {
  const logout = useCallback(async () => {
    try {
      const response = await fetch(`${process.env.REACT_APP_API_URL}/auth/logout`, { method: "POST", mode: "cors", credentials: 'include' });
      if (!response.ok) {
        throw new Error('Network response was not ok');
      }
    } catch (error) {
      console.log(error);
    }

  }, []);
  return logout;
}
