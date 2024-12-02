import { useState, useEffect } from 'react';

export interface NewBusinessInfo {
    name: string,
    desc: string,
    website: string
}

export function useCreateBusiness() {
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    // Function to send a POST request to create a new business
    const createBusiness = async (businessInfo: NewBusinessInfo) => {
        setIsLoading(true);
        setError(null);

        try {
            const response = await fetch("http://127.0.0.1:8080/users/0/businesses", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(businessInfo),
                mode: "cors",
                credentials: 'include',
            });

            if (!response.ok) {
                throw new Error('Failed to create business');
            }

            // Optionally handle the response data here
            const data = await response.json();
            console.log('Business created:', data);

            // Return the created business or any other data you want to handle
            return data;
        } catch (error: any) {
            setError(error.message || 'Something went wrong');
        } finally {
            setIsLoading(false);
        }
    };

    return { createBusiness, isLoading, error };
}