import { useState } from 'react';

export function useUploadImage() {
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    // Function to send a POST request to create a new business
    const uploadImage = async (image: any, business_id: string) => {
        setIsLoading(true);
        setError(null);

        const form_data = new FormData();
        form_data.append("image", image);

        try {
            const response = await fetch(`${process.env.REACT_APP_API_URL}/businesses/${business_id}/upload-image`, {
                method: "POST",
                body: form_data,
                mode: "cors",
                credentials: 'include',
            });

            if (!response.ok) {
                throw new Error('Failed to create business');
            }
        } catch (error: any) {
            setError(error.message || 'Something went wrong');
        } finally {
            setIsLoading(false);
        }
    };

    return { uploadImage, isLoading, error };
}
