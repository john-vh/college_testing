import { useState, useEffect, useMemo } from 'react';
import { Role } from '../components/InfoPage.tsx';

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

export function useIsAdmin(): Role {
    const accountInfo = useAccountInfo();
    return useMemo(() => accountInfo?.roles.includes("admin"), [accountInfo?.roles]) ? Role.Admin : Role.User;
}

export default useAccountInfo;
