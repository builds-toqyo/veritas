'use client';

import { motion } from 'framer-motion';
import { fadeIn } from '@/lib/animations';

export function StatsSection() {
  return (
    <section className='border-y border-muted bg-light py-20'>
      <div className='container max-w-6xl mx-auto px-4'>
        <div className='mx-auto grid gap-16 md:grid-cols-3 max-w-4xl'>
          <motion.div 
            initial='initial'
            whileInView='animate'
            viewport={{ once: true }}
            variants={fadeIn}
            className='text-center'
          >
            <h3 className='text-4xl font-bold text-primary'>99.9%</h3>
            <p className='mt-2 text-muted'>Uptime</p>
          </motion.div>
          <motion.div 
            initial='initial'
            whileInView='animate'
            viewport={{ once: true }}
            variants={fadeIn}
            className='text-center'
          >
            <h3 className='text-4xl font-bold text-primary'>24/7</h3>
            <p className='mt-2 text-muted'>Monitoring</p>
          </motion.div>
          <motion.div 
            initial='initial'
            whileInView='animate'
            viewport={{ once: true }}
            variants={fadeIn}
            className='text-center'
          >
            <h3 className='text-4xl font-bold text-primary'>ML v2.0</h3>
            <p className='mt-2 text-muted'>Latest Model</p>
          </motion.div>
        </div>
      </div>
    </section>
  );
}
