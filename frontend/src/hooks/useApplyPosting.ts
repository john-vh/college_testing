import { useState, useEffect } from 'react';
import { PostingApplicationInfo } from './useApplicationInfo';
import { PostingInfo } from './useAllPostings';

export function useApplyPosting() {

    const applyPosting = (post: PostingInfo, notes: string) => {
        const { business_id, id } = post;

        async function fetchData() {
            try {
                const response = await fetch(`http://127.0.0.1:8080/businesses/${business_id}/posts/${id}/apply`, {
                    method: 'POST',
                    mode: 'cors',
                    credentials: 'include',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({ notes })
                });
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
            } catch (error) {
                console.log(error);
            }
        }
        fetchData();
    };

    return applyPosting;
}
