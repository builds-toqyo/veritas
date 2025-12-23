'use client';

import Link from 'next/link';
import { motion } from 'framer-motion';
import { Button } from '@/components/ui/button';
import { fadeIn, stagger } from '@/lib/animations';

export function HeroSection() {
  return (
    <section className='hero-section relative min-h-[80vh] flex items-center justify-center'>
      <div className='absolute inset-0 bg-[url("/grid.svg")] opacity-10 bg-center [mask-image:linear-gradient(180deg,white,rgba(255,255,255,0))]' />
      <motion.div 
        initial='initial'
        animate='animate'
        variants={stagger}
        className='container max-w-6xl mx-auto px-4'
      >
        <div className='flex flex-col items-center gap-8 text-center mx-auto max-w-3xl'>
          <motion.h1 
            variants={fadeIn}
            className='text-4xl font-bold tracking-tighter sm:text-5xl md:text-6xl lg:text-7xl text-light'
          >
            Real-Time Risk Assessment for{' '}
            <span className='text-gradient'>
              Veritas Vault
            </span>
          </motion.h1>
          <motion.p 
            variants={fadeIn}
            className='mx-auto max-w-[700px] text-lg md:text-xl text-light opacity-80'
          >
            Advanced ML-powered risk analysis and monitoring system for DeFi operations.
          </motion.p>
          <div className='mt-6'>
            <motion.div variants={fadeIn}>
              <Button 
                asChild 
                size='lg' 
                className='btn-primary'
              >
                <Link href='/dashboard'>Launch Dashboard</Link>
              </Button>
            </motion.div>
          </div>
        </div>
      </motion.div>
    </section>
  );
}
