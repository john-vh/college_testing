import React, { useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { LandingNavbar } from '../components/LandingNavbar.tsx';
import { PostingContent } from '../components/PostingContent.tsx';
import { useLocation } from 'react-router-dom';
import { PostingInfo } from '../hooks/useAllPostings.ts';

export const Posting = () => {
    const location = useLocation();
    const navigate = useNavigate();
    const { post } = location.state as { post: PostingInfo } || {};

    useEffect(() => {
        if (!post) {
            navigate('/');
        }
    }, [post, navigate]);

    if (!post) { return; }

    return (
        <div>
            <LandingNavbar />
            <PostingContent post={post} />
        </div>
    );
};

export default Posting;
