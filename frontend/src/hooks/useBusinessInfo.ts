import { useState, useEffect, useMemo, useCallback } from 'react';

export interface BusinessInfo {
    id: string,
    user_id: string,
    created_at: string,
    status: string,
    name: string,
    desc: string,
    website: string
}

interface BusinessInfoProps {
    isAdmin?: boolean
}

export function useBusinessInfo({ isAdmin = false }: BusinessInfoProps) {
    const [businessInfo, setBusinessInfo] = useState<BusinessInfo[]>([]);

    const fetchData = useCallback(async () => {
        try {
            let response: Response | null = null;
            if (isAdmin) {
                response = await fetch("http://127.0.0.1:8080/admin/businesses", { mode: "cors", credentials: 'include' });
            }
            else {
                response = await fetch("http://127.0.0.1:8080/users/0/businesses", { mode: "cors", credentials: 'include' });
            }
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            const data: BusinessInfo[] = await response.json();
            setBusinessInfo(data);
        } catch (error) {
            console.log(error);
        }
    }, [isAdmin]);

    useState(() => {
        fetchData();
    });


    return { businessInfo, fetchData };
}

export function useIsFounder() {
    const { businessInfo } = useBusinessInfo({});

    // Memoize the check whether the list is empty or not
    const isFounder = useMemo(() => businessInfo.length > 0, [businessInfo]);

    return isFounder;
}

export function useUpdateBusiness() {
    const [loading, setLoading] = useState(false);

    const updateBusiness = (business: BusinessInfo) => {
        setLoading(true);
        const { id } = business;
        console.log(business);
        async function fetchData() {
            try {
                const response = await fetch(`http://127.0.0.1:8080/businesses/${id}`,
                    {
                        method: "PATCH", mode: "cors", credentials: 'include',
                        headers: {
                            "Content-Type": "application/json",
                        },
                        body: JSON.stringify(business),
                    });
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                else {
                    console.log(response);
                }
            } catch (error) {
                console.log(error);
            }
            finally {

                setLoading(false);
            }
        }
        fetchData();
    };

    return { updateBusiness, loading };
}