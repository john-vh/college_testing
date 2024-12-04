import { useState, useEffect } from 'react';
import { PostingApplicationInfo } from './useApplicationInfo';
import { PostingInfo } from './useAllPostings';
import { Toaster, Position, Intent } from "@blueprintjs/core";

const AppToaster = Toaster.create({
    position: Position.BOTTOM_RIGHT,
});

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
                AppToaster.show({
                    message: "Application submitted successfully!",
                    intent: Intent.SUCCESS,
                });
            } catch (error) {
                console.log(error);
                AppToaster.show({
                    message: "Failed to apply, application already submitted",
                    intent: Intent.DANGER,
                });
            }
        }
        fetchData();
    };

    return applyPosting;
}
