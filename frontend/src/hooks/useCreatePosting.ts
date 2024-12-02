import { useState, useEffect } from 'react';

export interface NewPostingInfo {
    title: string,
    desc: string,
    pay: number,
    time_est: number
    business_id: string
}

export function useCreatePosting() {
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    // Function to send a POST request to create a new business
    const createPosting = async (businessInfo: NewPostingInfo) => {
        setIsLoading(true);
        setError(null);

        try {
            const response = await fetch(`http://127.0.0.1:8080/businesses/${businessInfo.business_id}/posts`, {
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
            console.log('Posting created:', data);

            // Return the created business or any other data you want to handle
            return data;
        } catch (error: any) {
            setError(error.message || 'Something went wrong');
        } finally {
            setIsLoading(false);
        }
    };

    return { createPosting, isLoading, error };
}