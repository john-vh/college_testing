import { useState, useEffect } from 'react';

const useAccountInfo = () => {
    const [data, setData] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await fetch('http://127.0.0.1:8080/users', { mode: "cors", credentials: 'include' });
                if (!response.ok) {
                    console.log(response);
                    throw new Error('Network response was not ok');
                }
                const data = await response.json();
                setData(data);
            } catch (error) {
                setError(error);
            } finally {
                setLoading(false);
            }
        };

        fetchData();
    }, []);
    console.log(data);
    return { data, loading, error };
};

export default useAccountInfo;
