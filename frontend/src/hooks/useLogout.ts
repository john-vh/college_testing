import { useCallback } from 'react';

export function useLogout() {
    const logout = useCallback(async () => {
        try {
            const response = await fetch(`http://127.0.0.1:8080/auth/logout`, { method: "POST", mode: "cors", credentials: 'include' });
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
        } catch (error) {
            console.log(error);
        }

    }, []);
    return logout;
}