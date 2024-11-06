import React from 'react';
import { useParams } from 'react-router-dom';
import { LandingNavbar } from '../components/LandingNavbar.tsx';
import { PostingContent } from '../components/PostingContent';

export const Posting = () => {
    const { id } = useParams();

    return (
        <div>
            <LandingNavbar />
            <PostingContent id={id} />
        </div>
    );
};

export default Posting;
