import { useState, useEffect } from 'react';
import { PostingApplicationInfo } from './useApplicationInfo';

interface RejectApplicationProps {
    entry: PostingApplicationInfo,
    index: number
}

export function useRejectApplication() {

    const rejectApplication = (data: RejectApplicationProps) => {
        const { business_id, post_id } = data.entry;
        const { id: user_id } = data.entry.applications[data.index].user;

        async function fetchData() {
            try {
                const response = await fetch(`http://127.0.0.1:8080/businesses/${business_id}/posts/${post_id}/applications/${user_id}/reject`, { mode: "cors", credentials: 'include' });
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
            } catch (error) {
                console.log(error);
            }
        }
        fetchData();
    }; // Empty dependency array ensures this runs only once

    return rejectApplication;
}
