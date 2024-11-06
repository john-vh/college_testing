import { useState, useEffect } from 'react';

interface AccountInfo {
    id: string,
    created_at: string,
    status: number,
    email: string,
    name: string,
    email_verified: boolean,
    accounts: []
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

export default useAccountInfo;
