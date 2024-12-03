import { useState, useEffect, useMemo } from 'react';
import { Role } from '../components/InfoPage.tsx';
import { useIsFounder } from './useBusinessInfo.ts';

export interface AccountInfo {
    id: string,
    created_at: string,
    status: string,
    email: string,
    name: string,
    email_verified: boolean,
    roles: string[]
}

function useAccountInfo(): AccountInfo | null {
    const [accountInfo, setAccountInfo] = useState<AccountInfo | null>(null);

    useEffect(() => {
        async function fetchData() {
            try {
                const response = await fetch('http://127.0.0.1:8080/users', { mode: "cors", credentials: 'include' });
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                const data: AccountInfo = await response.json();
                setAccountInfo(data);
            } catch (error) {
                console.log(error);
            }
        }
        fetchData();
    }, []); // Empty dependency array ensures this runs only once

    return accountInfo;
}

export function useGetRole(): Role {
    const accountInfo = useAccountInfo();
    const isFounder = useIsFounder();

    return useMemo(() => {
        if (accountInfo?.roles.includes("admin")) {
            return Role.Admin; // Return Role.Admin if "admin" is in roles
        }
        if (isFounder) {
            return Role.Founder; // Return Role.Founder if businesses exist and no "admin" role
        }
        return Role.User; // Return Role.User otherwise
    }, [accountInfo?.roles, isFounder]);
}

export default useAccountInfo;
